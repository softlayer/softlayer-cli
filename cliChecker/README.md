## CLI Checker

The purpose of this little program is to compare the actual output of the commands/arguments of the CLI overtime.
It was introduced as part of the CLI refactor project, and might be removed later if its not really needed.

from the softlayer-cli root directory, do the following after making changes.

```
cd cliChecker
go build
./cliChecker > cliOutput.txt
git diff cliOutput.txt
```

Hopefully you should only see changes you expect.
