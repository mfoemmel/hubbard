project_name = read_project_name
repository_name = read_repository_name
authorize(project_name, 'admin')
dir = find_repository_dir(project_name, repository_name)
FileUtils.mkdir_p(dir)
Dir.chdir(dir) do
  exit $? unless system "git init --bare" 
  exec "git config hubbard.forkid #{project_name}/#{repository_name}/#{Time.now.to_i}" 
end
