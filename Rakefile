task :command_exists, [:command] do |_, args|
    abort "#{args.command} doesn't exists" if `command -v #{args.command} > /dev/null 2>&1 && echo $?`.chomp.empty?
end

task :has_mockery do
    Rake::Task['command_exists'].invoke('mockery')
end

task :has_gsed do
    Rake::Task['command_exists'].invoke('gsed')

end
namespace :mock do
    namespace :mockery do
      desc "upgrade mockery version"
      task :upgrade => [:has_mockery] do
        system %{
          brew update && 
          brew upgrade mockery && 
          brew cleanup mockery
        }
      end
      
      desc "show installed version"
      task :show_version => [:has_mockery] do
        system "mockery --version --log-level error"
      end
    end
  
    namespace :generate do
      desc "generate mock for all"
      task :all => [:has_mockery] do
        rm_rf %w(mocks)
        system %{
          mockery --output ./mocks --case=underscore --dir ./src/migrationtool/ --all --keeptree --recursive && \
          mockery --output ./mocks/db --case=underscore --dir ./src/db/ --all --keeptree --recursive
        }
      end
    end
end

desc "run tests"
task :test => [:has_gsed] do
  system %{
    color_red=$'\e[0;31m'
    color_yellow=$'\e[0;33m'
    color_white=$'\e[0;37m'
    color_off=$'\e[0m'
    
    any_errors="0"
    
    for s in $(go list ./...); do 
      if ! go test -failfast -p 1 -v -race "${s}"; then
        echo "\t\n${color_red}${s}${color_off} ${color_yellow}fails${color_off}...\n"
        any_errors="1"
        break
      fi
    done
    
    if [[ "${any_errors}" == "0" ]]; then
      echo "\n\n${color_white}Tests are passing...${color_off}"
      echo "${color_white}Calculating code coverage${color_off}"
      go test ./... -coverpkg=./src/... -coverprofile ./coverage.out > /dev/null 2>&1
      code_coverage_ratio=$(go tool cover -func ./coverage.out | grep "total:" | awk '{print $3}')
      echo "${color_white}Total test coverage: ${color_yellow}${code_coverage_ratio}${color_off}"
      code_coverage_ratio_md=${code_coverage_ratio/%/25}
      gsed -i -r "s/coverage-[0-9\.\%]+/coverage-${code_coverage_ratio_md}/" README.md &&
      echo "README updated...\n"
    fi
  }
end