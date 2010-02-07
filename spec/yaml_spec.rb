require 'spec_helper'
require 'yaml'

YAML_OPTION=" -f yaml "
describe "Hubble with yaml output" do
  before(:each) do
    reset_file_system
  end

  after(:all) do
    reset_file_system
  end

  it "should list-projects" do
    hub("yammer", "create-project a a-desc")
    hub("yammer", "create-project b b-desc")
    hub("yammer", "create-project c c-desc")

    projects = YAML::load(hub("yammer", "#{YAML_OPTION} list-projects")).map{|project|project[:name]}
    projects.should == ["a", "b", "c"]
  end

  it "should list repositories" do
    hub("yammer", "create-project a a-desc")
    hub("yammer", "create-repository a b")

    repositories = YAML::load(hub("yammer", "#{YAML_OPTION} list-repositories a"))
    repositories.length.should == 1
    repositories.first[:name].should == "b"
    repositories.first[:url].should == "#{ENV['USER']}@#{HUB_HOST}:a/b.git"
  end

  it "should list permissions" do
    hub("yammer", "create-project a a-desc")
    permissions = YAML::load(hub("yammer", "#{YAML_OPTION} list-permissions a"))
    permissions.length.should == 1
    permissions.first[:user].should == "yammer"
    permissions.first[:access].should == "admin"
  end

  it "should list users for admin" do
    hub("kipper", "add-key laptop", "ssh-rsa yabbadabba fdsa")
    hub("yammer", "add-key ipad", "ssh-rsa fadadada chacaha")
    users = YAML::load(hub("admin", "#{YAML_OPTION} list-users"))
    users.size.should == 2
    users.should == ["kipper", "yammer"]
  end
end
