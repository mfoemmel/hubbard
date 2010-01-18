project_name = read_project_name
authorize(project_name, 'admin')
dir = find_project_dir(project_name)
username = ARGV.shift
action = ARGV.shift
unless ['admin','write','read'].member?(action)
  $stderr.puts "Not a valid action (must be one of: read, write, admin)"
  exit 1
end
File.open(File.join(dir, ".permissions"), "r+") do |f| 
  f.flock(File::LOCK_EX)
  contents = f.read
  puts contents
  f.flock(File::LOCK_UN)
end
