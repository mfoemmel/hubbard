from_project_name = read_project_name
from_repository_name = read_repository_name
to_project_name = read_project_name
to_repository_name = read_repository_name
authorize(from_project_name, 'read')
authorize(to_project_name, 'admin')
from_dir = find_repository_dir(from_project_name, from_repository_name)
to_dir = find_repository_dir(to_project_name, to_repository_name)
forkid = Dir.chdir(from_dir) { `git config --get hubbard.forkid` }
FileUtils.mkdir_p(to_dir)
exit $? unless system "git clone --bare --shared #{from_dir} #{to_dir}" 
Dir.chdir(to_dir) do
  exec "git config hubbard.forkid #{forkid}"
end

