project_name = read_project_name
authorize(project_name, 'admin')
dir = find_project_dir(project_name)
username = ARGV.shift

permissions_file = File.join(dir, ".permissions")
if File.exists?(permissions_file)
  File.open(permissions_file, "r+") do |f|
    f.flock(File::LOCK_EX)
    contents = f.read
    puts contents
    f.flock(File::LOCK_UN)
  end
end
