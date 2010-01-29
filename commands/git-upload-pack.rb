unless ARGV.shift =~ /(#{Hubbard::PROJECT_REGEX})\/(#{Hubbard::REPOSITORY_REGEX}).git/
  $stderr.puts "Repository not found"
  exit 1
end
project_name = $1
repository_name = $2
authorize(project_name, 'read')
dir = find_repository_dir(project_name, repository_name)
exec "git-upload-pack #{dir}"
