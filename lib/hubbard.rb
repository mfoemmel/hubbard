require 'socket'

module Hubbard
  PROJECT_REGEX='[a-zA-Z0-9\-]{1,32}'
  REPOSITORY_REGEX='[a-zA-Z0-9\-]{1,32}'
  USERNAME_REGEX='[a-zA-Z0-9\-]{1,32}'
  KEY_NAME_REGEX = /[a-zA-Z0-9]+/
  KEY_REGEX = /(ssh-rsa|ssh-dsa) ([a-zA-Z0-9\+\/]+[=]*)/

  HUB_DATA = ENV['HUB_DATA'] || File.expand_path("~/.hubbard")
  HUB_HOST = ENV['HUB_HOST'] || Socket.gethostname
  PROJECTS_PATH = File.join(HUB_DATA, "projects")
  ACCOUNTS_PATH = File.join(HUB_DATA, "accounts")

  ACTIONS = ['read', 'write', 'admin']
end
