require 'fileutils'

project_name = read_project_name
description = next_arg "Please specify a project description"
dir = find_project_dir(project_name)
if File.exist?(dir)
  $stderr.puts "Project already exists with that name"
  exit 4
end
unless Dir.mkdir(dir)
  $stderr.puts "Unable to create directory: #{dir}"
end
visibility = OPTIONS[:private] ? "private" : "public"
if @username != 'admin'
  File.open(File.join(dir, ".permissions"), "w") { |f| f << "#{@username}=admin\n" }
end
File.open(File.join(dir, ".visibility"), "w") { |f| f << "#{visibility}\n" }
File.open(File.join(dir, ".description"), "w") { |f| f << "#{description}\n" }
