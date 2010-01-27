project_name = read_project_name
authorize(project_name, 'read')
Dir.foreach(find_project_dir(project_name)) do |repository_name|
  next if repository_name =~ /^\./
  git_url = "#{ENV['USER']}@#{HUB_HOST}:#{project_name}/#{repository_name}.git"
  puts "#{repository_name}\t#{git_url}"
end
