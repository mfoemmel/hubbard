require 'fileutils'

project_name = read_project_name
authorize(project_name, 'admin')

description = next_arg "Please specify the description"

dir = find_project_dir(project_name)
if !File.exist?(dir)
  $stderr.puts "Project not found"
  exit 4
end

File.open(File.join(dir, ".description"), "w") { |f| f << description << "\n" }

