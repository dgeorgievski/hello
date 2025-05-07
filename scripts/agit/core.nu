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
