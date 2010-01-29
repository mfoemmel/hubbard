require 'socket'

module Hubbard
  PROJECT_REGEX='[a-zA-Z0-9\-]{1,32}'
  REPOSITORY_REGEX='[a-zA-Z0-9\-]{1,32}'
  USERNAME_REGEX='[a-zA-Z0-9\-]{1,32}'

  HUB_DATA = ENV['HUB_DATA'] || File.expand_path("~/.hubbard")
  HUB_HOST = ENV['HUB_HOST'] || Socket.gethostname
end
