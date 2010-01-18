unless ARGV.shift =~ /(#{PROJECT_REGEX})\/(#{REPOSITORY_REGEX}).git/
  $stderr.puts "Repository not found"
  exit 1
end
project_name = $1
repository_name = $2
authorize(project_name, 'write')
dir = find_repository_dir(project_name, repository_name)
exec "git-receive-pack #{dir}"
