from_project_name = read_project_name
from_repository_name = read_repository_name

to_project_name = read_project_name
to_repository_name = read_repository_name

from_dir = find_repository_dir(from_project_name, from_repository_name)
to_dir = find_repository_dir(to_project_name, to_repository_name)

authorize(from_project_name, 'admin')
authorize(to_project_name, 'admin')

if not File.exist?(from_dir)
  error 4, "Repository not found"
end

if File.exist?(to_dir)
  error 4, "Repository already exists with that name"
end

FileUtils.mv(from_dir, to_dir)
