require 'spec_helper'
require 'yaml'

YAML_OPTION=" -f yaml "
describe "Hubble with yaml output" do
  before(:each) do
    reset_file_system
  end

  it "should load list-projects" do
    hub("yammer", "create-project a a-desc")
    hub("yammer", "create-project b b-desc")
    hub("yammer", "create-project c c-desc")

    projects = YAML::load(hub("yammer", "#{YAML_OPTION} list-projects")).map{|project|project[:name]}
    projects.should == ["a", "b", "c"]
  end

  it "should create repositories" do
    hub("yammer", "create-project a a-desc")
    hub("yammer", "create-repository a b")

    repositories = YAML::load(hub("yammer", "#{YAML_OPTION} list-repositories a"))
    repositories.length.should == 1
    repositories.first[:name].should == "b"
    repositories.first[:url].should == "#{ENV['USER']}@#{HUB_HOST}:a/b.git"
  end
end
