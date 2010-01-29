require 'spec_helper'
require 'yaml'

YAML_OPTION=" -f yaml "
describe "Hubble with yaml output" do
  before(:each) do
    reset_file_system
  end

  it "should load list-projects" do
    hub("yammer", "create-project a")
    hub("yammer", "create-project b")
    hub("yammer", "create-project c")

    projects = YAML::load(hub("#{YAML_OPTION} yammer", "list-projects"))
    projects.should == ["a", "b", "c"]
  end

  it "should create repositories" do
    hub("yammer", "create-project a")
    hub("yammer", "create-repository a b")

    repositories = YAML::load(hub("yammer", "#{YAML_OPTION} list-repositories a"))
    repositories.length.should == 1
    repositories.first[:name].should == "b"
    repositories.first[:url].should == "#{ENV['USER']}@#{HUB_HOST}:a/b.git"
  end
end
