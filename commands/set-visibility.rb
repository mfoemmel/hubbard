require 'fileutils'

project_name = read_project_name
visibility = next_arg "Please specify one of: public, private"
if visibility != 'public' and visibility != 'private'
  error 3, "Please specify one of: public, private"
end

dir = find_project_dir(project_name)
if !File.exist?(dir)
  error 4, "Project not found"
end

File.open(File.join(dir, ".visibility"), "w") { |f| f << visibility << "\n" }

