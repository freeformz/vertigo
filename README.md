vertigo
=======

Just another blog engine in Go.

Currently under heavy development!

Vertigo aims to bring simpler way to publish simple websites or blogs, which are normally done with heavy content management systems like Wordpress. Theme development should be distinct from backend code so that any knowledge in Go would not be needed when developing themes.

Once more stable, Vertigo is supposed to be available in simple executable file for platforms supported by Go. Any unnecessary 3rd party packages are subject to be removed to further increase easy portability.

The current database choice is MongoDB because of changing database structures. I'm yet to choose whether that will change once the structures are complete.

Vertigo is supposed to have an frontend API so it can be hooked to various frontend JavaScript frameworks. However, all functionality should also be available without one.


##What works

Since there is so much still to do, I'd rather list some things that are currently working:

- User CRUD options
- Basic session control
- Creating new posts, listing them on each user page


##Known bugs

There are some known bugs or other quirks which would need resolving:

- ~~In Martini, how you require field but do not insert it into database - does it require seperate structures? Currently saving a user will leave a empty field password field into database~~
- In Martini, how you render multiple layouts without calling r.HTML multiple times?


##License

MIT