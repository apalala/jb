
1. I noted that agents would use *nix pipelines that included commands like sed to edit files instead of using the tools provided by the editor, so I forbid that in RULES.md.
1. The bots would constantly change files outside the scope of the current task so I resourted to `chmod -r ag-w notinscope`. The bots would run `chmod` to access the files.
1. When a bot made a mistake it resorted to a blind `git checkout` to roll back its changes effectively destroying any work in progress different from their incorrect changes.
1. After adding rules regarding misbehavior to RULES.md, the bots would continue to step outside of the stated boundaries. When asked "Why?" they would apologize and promise they would not do it again, but they would.
1. I started placing stubs on `PATH` for the forbidden commands, but:

    a. That limited legitimate use, like `git status`.
    b. The bots would use `which` to find other versions of the commands, and would try to examine the contents of `~/bin/cmd` to bypass the stubs.
1. Finally I ditched the stubs in favor for smart shells that will let legitimate uses pass and will output several seconds of literary works randomized with a blue noise generator for rouge behavior.
1. I called the project "Johannes Blues" after Johannes Gutemberg, and carefully choose the works that are streamed to the bots. The implicit joke is that I'm helping with their literary culture.
1. The reaction from the bots has been mostly to try to bypass a few times, and then just comply with the rules of the IDE and documents like RULES.md.
1. Whe I find a shell is restricting legitimate use or just being obnoxious, I easy the restrictions right away.
