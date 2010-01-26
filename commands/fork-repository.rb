from_project_name = read_project_name
from_repository_name = read_repository_name
to_project_name = read_project_name
to_repository_name = read_repository_name
authorize(from_project_name, 'read')
authorize(to_project_name, 'admin')
from_dir = find_repository_dir(from_project_name, from_repository_name)
to_dir = find_repository_dir(to_project_name, to_repository_name)
FileUtils.mkdir_p(to_dir)
exec "git clone --bare #{from_dir} #{to_dir}" 

