require 'rubygems'
require 'rake'

begin
  require 'jeweler'
  Jeweler::Tasks.new do |gem|
    gem.name = "hubbard"
    gem.summary = %Q{Hubbard is a command line tool for managing git repositories.}
    gem.description = %Q{Hubbard is a command line tool for managing git repositories.}
    gem.email = "git@foemmel.com"
    gem.homepage = "http://github.com/mfoemmel/hubbard"
    gem.authors = ["Matthew Foemmel"]
    gem.add_development_dependency "rspec", ">= 1.2.9"
    gem.files << Dir['commands/*.rb']
    gem.bindir = 'bin'
    gem.executables << 'hubbard'
    gem.require_path = ''
  end
  Jeweler::GemcutterTasks.new
rescue LoadError
  puts "Jeweler (or a dependency) not available. Install it with: gem install jeweler"
end

require 'spec/rake/spectask'
Spec::Rake::SpecTask.new(:spec) do |spec|
  spec.libs << 'lib' << 'spec'
  spec.spec_files = FileList['spec/**/*_spec.rb']
end

Spec::Rake::SpecTask.new(:rcov) do |spec|
  spec.libs << 'lib' << 'spec'
  spec.pattern = 'spec/**/*_spec.rb'
  spec.rcov = true
end

task :spec => :check_dependencies

task :default => :spec

require 'rake/rdoctask'
Rake::RDocTask.new do |rdoc|
  version = File.exist?('VERSION') ? File.read('VERSION') : ""

  rdoc.rdoc_dir = 'rdoc'
  rdoc.title = "hubbard #{version}"
  rdoc.rdoc_files.include('README*')
  rdoc.rdoc_files.include('lib/**/*.rb')
end
