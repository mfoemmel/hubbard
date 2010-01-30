require 'fileutils'

project_name = read_project_name
visibility = next_arg "Please specify one of: public, private"
if visibility != 'public' and visibility != 'private'
  $stderr.puts "Please specify one of: public, private"
  exit 3
end

dir = find_project_dir(project_name)
if !File.exist?(dir)
  $stderr.puts "Project not found"
  exit 4
end

File.open(File.join(dir, ".visibility"), "w") { |f| f << visibility << "\n" }

