use common.nu *
use complete.nu *
use stat.nu *

# Check if last commit message complies with conventional commits specs
export def is_conventional_commit [] {
  let subject = git show -s --format="%s"

  # Body is optional
  # let body = git show -s --format=%b

  # Footer with trailers is optional
  # Signed-off-by: Alice <alice@example.com>
  
  # check if subject complies
  # Examples https://www.conventionalcommits.org/en/v1.0.0/#examples
  #   returns true:
  #     fix(api)!: send an email to the customer when a product is shipped"
  #     chore!: drop support for Node 6
  #     feat: allow provided config object to extend other configs
  #   returns false:
  #     another update
  $subject | parse --regex '(?P<type>^[a-z]+)(?P<scope>\(.*\))?(?P<break>!?):(?P<desc>.*$)' | is-not-empty
}

export def is_conventional_commit_hook [commit_msg_file] {
  let commit_msg = open $commit_msg_file | lines | first 

  # Body is optional
  
  # Footer with trailers is optional
  # Signed-off-by: Alice <alice@example.com>
  
  # check if subject complies
  # Examples https://www.conventionalcommits.org/en/v1.0.0/#examples
  #   returns true:
  #     fix(api)!: send an email to the customer when a product is shipped"
  #     chore!: drop support for Node 6
  #     feat: allow provided config object to extend other configs
  #   returns false:
  #     another update
  $commit_msg | parse --regex '(?P<type>^[a-z]+)(?P<scope>\(.*\))?(?P<break>!?):(?P<desc>.*$)' | is-not-empty
}

export def git-squash-last [
    num:int
] {
    let l = git log  --pretty=»»¦««%s»¦«%b -n $num
    | split row '»»¦««' | slice 1..
    | split column '»¦«' message body
    | each { $"($in.message)\n\n($in.body)" }
    | str join "\n===\n"
    git reset --soft $"HEAD~($num)"
    git commit --edit -m $l
}

export def git-merge [
    branch?:            string@cmpl-git-branches
    --abort (-a)
    --continue (-c)
    --quit (-q)
    --squash (-s)
    --fast-farward (-f)
    --remote (-r)='origin':  string@cmpl-git-remotes
] {
    mut args = []
    if $squash { $args ++= [--squash] }
    if $fast_farward { $args ++= [--ff] } else { $args ++= [--no-ff] }
    if ($branch | is-empty) {
        git merge ...$args $"($remote)/(git_main_branch)"
    } else {
        git merge ...$args $branch
    }
    if $squash {
        git commit -v
    }
}
