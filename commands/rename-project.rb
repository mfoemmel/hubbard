from_project_name = read_project_name
to_project_name = next_arg "Please specify the new project name"
from_dir = find_project_dir(from_project_name)
to_dir = find_project_dir(to_project_name)
if File.exist?(to_dir)
  error 3, "A project already exists with that name"
end

authorize(from_project_name, 'admin')
FileUtils.mv(from_dir, to_dir)
