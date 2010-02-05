$LOAD_PATH.unshift(File.dirname(__FILE__))
$LOAD_PATH.unshift(File.join(File.dirname(__FILE__), '..', 'lib'))
require 'hubbard'
require 'spec'
require 'spec/autorun'

Spec::Runner.configure do |config|
end

# Something in the Rakefile is generating
# a bunch of GIT_* environment variables,
# which mess everything up, so undo.
ENV.each do |key,value|
  ENV[key] = nil if key =~ /^GIT_/
end

HUB=File.expand_path(File.join(File.dirname(__FILE__), "..", "bin", "hubbard"))
ENV['HUB_USER'] = HUB_USER = "hub"
ENV['HUB_HOST'] = HUB_HOST = "example.com"
ENV['HUB_DATA'] = HUB_DATA =  File.expand_path(File.join(File.dirname(__FILE__), '..', "data"))
ENV['GIT_SSH'] = File.expand_path(File.join(File.dirname(__FILE__), "gitssh"))

def hub(username, command, input=nil)
  if input
    result = `echo #{input} | #{HUB} #{username} #{command}`
  else
    result = `#{HUB} #{username} #{command}`
  end
  handle_exitstatus(command, result)
  result
end

def git(username, command)
  ENV['HUB_USERNAME'] = username
  result = `git #{command}`
  handle_exitstatus(command, result)
  result
end

def reset_file_system
  FileUtils.rm_rf HUB_DATA
  FileUtils.rm_rf "tmp"
end

def list_projects(user)
  hub(user, "list-projects").split("\n").map { |line| line.split[0] }
end

private
def handle_exitstatus(command, result)
  exit_code = $?.exitstatus
  if exit_code != 0
    msg = "Command failed [#{exit_code}]: git #{command}:\n#{result}"
    raise SystemCallError.new(msg, exit_code)
  end
end
