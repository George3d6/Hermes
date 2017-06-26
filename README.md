<img src="client/static/logo_mini.png" alt="" style="height:220px !important; display:inline-block">

#Hermes

A minimalistic, easy to deploy, easy to use, self-contained file server.

#####Version 0.4.2

##Instalation

```bash
wget 'https://hermes.cerebralab.com/get/file/?file=Hermes.tar.gz' && tar -xvf Hermes.tar.gz && cd dist && ./hermes 'config.json' 'admin' 'password'
```

That's it, you've deployed the Hermes file server, the last two parameters are the username and password of the initial administrator user, so consider changing them to something other than 'admin' and 'password'. To access it go to:

```bash
http://localhost:3280/get/authentication/?identifier=admin&credentials=password
```

This link should authenticate you and display the main web interface. There are 3 buttons on top, one is for signing in (you already did that by virtue of using the link before), one is for listing all your files(you have none at the moment) and one is for adding users.

Keep in mind if you are using a version prior to 1.0.0 you are likely to run into bugs or lack of a feature which you'd consider 'essential'. If any of the former happen to you please open an issue here and I will look into it as fast as possible.

##Roadmap

- Improve error&success communication to the client and display them in a friendly way to the user

- Add the option to allow tokens to access individual files

- Add loading screens

- Add an option to explore the various tokens to the UI

- Allow asynchronous upload for small files on all browsers (even those that don't fully support the File web API)

##Development

This is a rather small project and I am uncertain as to how easy it is to read&understand&expand, that being said if anyone feels bored or finds himself using this software a lot and needing some extra feature, I would be overjoyed to welcome them as a contributor. Hopefully the information found here will help with that.

###How to compile and run
####Dependencies

Currently there are the following dev dependencies in order to properly compile the server and client, in addition to this certain npm packages might have their own dependencies or restrictions that I am unaware of.

- node version 8.1.2 (Though any version starting with 6 should work)
- npm version 5.0.3 (Though any version whatsoever should work)
- go version 1.8.3 (Though any version starting with 1.0.0 should work)
- sassc version (?)
- A linux distro

####Server
The whole process is rather easy, simply build the .go files inside server and then run using the default config.json after making sure you have all the directories needed. The following commands should accomplish that:

```bash
git clone https://git.cerebralab.com/george/Hermes;
cd Hermes/server;
go build && mv server ../hermes;
cd ../;
mkdir ser storage;
./hermes 'config.json' 'admin' 'admin';
```

An easy way to test if the project is running on your machine is runing the regression

```bash
cd tests && npm install && ./regression.sh
```

####The client
Currently the js client is built using typescript and bundled using webpack, the css is compiled from scss, building all those components will include the following steps:
```bash
npm install; #Install all dependencies
webpack; #Compile typescript in the /ts folder into a bundle.js file placed in static
sassc scss/main.scss static/bundle.css; #Compile the scss into a bundle.css files placed in static
```

####Feedback and help
In theory this should also work on BSD&OSX without much additional work, if anyone can confirm that and/or tell me how I could make some of the linux specific things compatible and/or make a pull request in order to do so himself, I would very much welcome that.

As far as Windows is concerned I have no idea how much work would be needed to adapt this project to Windows since I haven't really got much experience coding on Windows, if you wish to make a pull request that makes this project windows-dev compatible I would also welcome that with open arms (though I'd have a hard time confirming if it works properly).

###Where is dev work required ?

####Testing

Currently this projects boast an incomplete regression and some defunct functional tests. I would like to increase the coverage of the regression and possibly split it into smaller tests.

Once the code stabilizes a bit and if this project pick up traction past one or two users some unit tests would also be very welcome

####Standardization and boilerplate removal

Currently the authentication layer and the file list layer use slightly different mechanisms in order to ensure thread safety, I'm considering adapting them to use a more similar style.

There's a lot of boilerplate right now, in part due to the fact that I dislike reducing code to self contained functions & routines before I have insighed in the most light-weight way to do that, with this kind of projects, I find, that require developing with some boilerplate first. I am working at removing it and any pull requests that remove boilerplate in good ways, that is, reduce the amount of code and improve readability and logic, will be very welcome.

Currently the front-end is a mess, it should be re-organized in to better named files and maybe modularized a bit more.

####Dev dependencies removal

Errors are handled using a personal logging & error handling library which, after working a bit on this project, I've realized is quite poorly written. For the sake of making this project more accessible I will work on taking that out and replacing it with boilerplate error handling or a more widely adopted library.

If possible removing the dependency on semantic-ui would be nice, since that drags along JQuery with it and I find it can slow load and processing times down quite significantly on mobile devices or in areas with terrible connection.

####Features

Basically the ones mentioned in the 'chapter' above, as to their implementation... that may vary.

The hardest feature to think about right now is the ability to dispense quick login url (e.g.: my.hermes.com/access/now/?id=myfirend&safety=aHash) without compromizing to much on the relative safety of the current auth system

There are also some deeply neste for loops in there that I'd love to get rid of, granted that it will not take away any generality (ie. force to write more boilerplate), after all even itterating through arrays of thousands of elements should take an
insignificant amount of time, golang arrays and maps are quite fast.

###Design goals

This file server was designed with two goals in mind:

1. Self containment, which means there should be no dependencies past compile time
2. Ease of deployment and migration

By virtue of the language it is written in it also happens to be reasonably fast and that is an accidental characteristic I would like to
maintain and improve.

Those are the main things you should keep in mind when doing pull requests.
