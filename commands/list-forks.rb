project_name = read_project_name
repository_name = read_repository_name
authorize(project_name, 'read')
forkid = Dir.chdir(find_repository_dir(project_name, repository_name)) { `git config --get hubbard.forkid` }
project_dir = find_project_dir(project_name)
Dir.foreach(File.join(HUB_DATA, 'projects')) do |dir|
  next if dir == "." || dir == ".."
  next unless is_authorized(dir, 'read')
  Dir.foreach(find_project_dir(project_name)) do |repository_name|
    next if repository_name =~ /^\./
    Dir.chdir(find_repository_dir(project_name, repository_name)) do
      if forkid == `git config --get hubbard.forkid`
        puts "#{project_name}/#{repository_name}"
      end
    end
  end
end
