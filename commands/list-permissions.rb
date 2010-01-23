project_name = read_project_name
authorize(project_name, 'admin')
dir = find_project_dir(project_name)
username = ARGV.shift
File.open(File.join(dir, ".permissions"), "r+") do |f| 
  f.flock(File::LOCK_EX)
  contents = f.read
  puts contents
  f.flock(File::LOCK_UN)
end
