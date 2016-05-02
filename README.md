gomoku project in golang
graphic interface in sdl

// Speeds up compilation
go install -v github.com/veandco/go-sdl2/{sdl_ttf,sdl}

// Sets up hook that runs unit tests
ln -s ../../pre-commit.sh .git/hooks/pre-commit
