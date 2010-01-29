require 'spec_helper'
require 'fileutils'

describe "Hubble" do
  before(:each) do
    reset_file_system
  end

  it "should create project" do
    hub("kipper", "create-project foo")
    projects = hub("kipper", "list-projects").split("\n")
    projects.should == ["foo"]
  end

  it "should not allow multiple projects with same name" do
    hub("kipper", "create-project foo")
    lambda { hub("kipper", "create-project foo") }.should raise_error
  end

  it "should delete project" do
    hub("kipper", "create-project foo")
    hub("kipper", "delete-project foo")

    projects = hub("kipper", "list-projects").split("\n")
    projects.should == []
  end

  it "should default to public project" do
    hub("kipper", "create-project foo")

    # Other users can see...
    projects = hub("tiger", "list-projects").split("\n")
    projects.should == ["foo"]

    # But not delete
    lambda { hub("tiger", "delete-project foo") }.should raise_error
  end

  it "should support private project" do
    hub("kipper", "create-project foo --private")

    # Other users can't see
    projects = hub("tiger", "list-projects").split("\n")
    projects.should == []
  end

  it "should create repositories" do
    hub("kipper", "create-project foo")
    hub("kipper", "create-repository foo bar")

    repositories = hub("kipper", "list-repositories foo").split("\n")
    repositories.length.should == 1
    name,url = repositories[0].split
    name.should == "bar"
    url.should == "#{ENV['USER']}@#{HUB_HOST}:foo/bar.git"    
  end

  def with_test_project
    Dir.mkdir('tmp')
    Dir.chdir('tmp') do
      File.open("README", "w") { |f| f << "Hello, world\n" }
      fail unless system "git init"
      fail unless system "git add README"
      fail unless system "git commit -m 'initial commit'" 
      yield
    end
  end

  it "should allow git push" do
    hub("kipper", "create-project foo")
    hub("kipper", "create-repository foo bar")

    with_test_project do
      git("kipper", "push #{ENV['USER']}@#{HUB_HOST}:foo/bar.git master")
    end
  end

  it "should allow git push with write permissions" do
    hub("kipper", "create-project foo")
    hub("kipper", "add-permission foo tiger write")
    hub("kipper", "create-repository foo bar")

    with_test_project do
      git("tiger", "push #{ENV['USER']}@#{HUB_HOST}:foo/bar.git master")
    end
  end

  it "should not allow git push with read permissions" do
    hub("kipper", "create-project foo")
    hub("kipper", "add-permission foo tiger read")
    hub("kipper", "create-repository foo bar")

    with_test_project do
      lambda { git("tiger", "push #{ENV['USER']}@#{HUB_HOST}:foo/bar.git master") }.should raise_error
    end
  end

  it "should allow git pull" do
    hub("kipper", "create-project foo")
    hub("kipper", "create-repository foo bar")

    with_test_project do
      git("kipper", "push #{ENV['USER']}@#{HUB_HOST}:foo/bar.git master")
      git("kipper", "pull #{ENV['USER']}@#{HUB_HOST}:foo/bar.git master")
    end
  end

  it "should not allow git pull with no permissions" do
    hub("kipper", "create-project foo --private")
    hub("kipper", "create-repository foo bar")

    with_test_project do
      git("kipper", "push #{ENV['USER']}@#{HUB_HOST}:foo/bar.git master")
      lambda { git("tiger", "pull #{ENV['USER']}@#{HUB_HOST}:foo/bar.git master") }.should raise_error
    end
  end

  it "should allow git pull with read permissions" do
    hub("kipper", "create-project foo")
    hub("kipper", "create-repository foo bar")

    with_test_project do
      git("kipper", "push #{ENV['USER']}@#{HUB_HOST}:foo/bar.git master")
      git("tiger", "pull #{ENV['USER']}@#{HUB_HOST}:foo/bar.git master")
    end
  end

  it "should fork repository in same project" do
    hub("kipper", "create-project foo")
    hub("kipper", "create-repository foo bar")

    with_test_project do
      git("kipper", "push #{ENV['USER']}@#{HUB_HOST}:foo/bar.git master")
      hub("kipper", "fork-repository foo bar foo bar2")
      git("kipper", "pull #{ENV['USER']}@#{HUB_HOST}:foo/bar2.git master")
    end
  end

  it "should fork repository in different project" do
    hub("kipper", "create-project foo")
    hub("kipper", "create-project foo2")
    hub("kipper", "create-repository foo bar")

    with_test_project do
      git("kipper", "push #{ENV['USER']}@#{HUB_HOST}:foo/bar.git master")
      hub("kipper", "fork-repository foo bar foo2 bar2")
      git("kipper", "pull #{ENV['USER']}@#{HUB_HOST}:foo2/bar2.git master")
    end
  end

  it "should track projects related by forking" do
    hub("kipper", "create-project foo")
    hub("kipper", "create-repository foo bar")

    with_test_project do
      git("kipper", "push #{ENV['USER']}@#{HUB_HOST}:foo/bar.git master")
      hub("kipper", "fork-repository foo bar foo bar2")
      hub("kipper", "list-forks foo bar").should == "foo/bar\nfoo/bar2\n"
    end
  end

  it "should require read access to fork repository" do
    hub("kipper", "create-project foo")
    hub("kipper", "create-project foo2")
    hub("kipper", "create-repository foo bar")

    with_test_project do
      git("kipper", "push #{ENV['USER']}@#{HUB_HOST}:foo/bar.git master")
      lambda { hub("tiger", "fork-repository foo bar foo2 bar2") }.should raise_error
      hub("kipper", "add-permission foo tiger read")
      lambda { hub("tiger", "fork-repository foo bar foo2 bar2") }.should raise_error
      hub("kipper", "add-permission foo2 tiger write")
      lambda { hub("tiger", "fork-repository foo bar foo2 bar2") }.should raise_error
      hub("kipper", "add-permission foo2 tiger admin")
      hub("tiger", "fork-repository foo bar foo2 bar2")
      hub("kipper", "add-permission foo2 tiger admin")      
    end
  end

  it "should remove permission" do
    hub("kipper", "create-project foo")
    hub("kipper", "create-repository foo bar")
    hub("kipper", "add-permission foo tiger read")
    hub("kipper", "remove-permission foo tiger")

    with_test_project do
      lambda { git("tiger", "push #{ENV['USER']}@#{HUB_HOST}:foo/bar.git master") }.should raise_error
    end
  end

  it "should add ssh key" do
    hub("kipper", "add-key laptop", "ssh-rsa yabbadabba fdsa")
  end

  it "should allow admin to run-as another user" do
    hub("admin", "run-as kipper create-project foo")
    projects = hub("kipper", "list-projects").split("\n")
    projects.should == ["foo"]
  end
end
