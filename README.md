# A fork of Navidrome Music Server
## An experiment with a different way to create your own smart playlists

### I endorse Navidrome!

- _Navidrome is an excellent open source web-based music collection server_
- _If you have a large music collection and want to host your own Subsonic compatible server, I recommend Navidrome!_
- _I have tried several self-hosted music servers and Navidrome is IMHO the best Subsonic compatible server._
- _It is also super easy to set up - I run mine in a Docker container on my NAS_

_You don't have to stick with the Navidrome Web interface either:_

In your web brwoser, [Airsonic (refix) UI](https://github.com/tamland/airsonic-refix) is a great open-source project that works seamlessly with Navidrome

And [Amperfy](https://github.com/BLeeEZ/amperfy) is a great open-source iOS client for iPhone and/or iPad, and again it works seamlessly with Navidrome

### A Brief History Of Playlists

- Basic playlists let you create a sequence of songs (tracks)...nice but how could we make them better?  Shuffle the order? Meh!
- Streaming music services have their own playlists and suggestions...some people love them but I do not find the "algorithm" all that helpful
- Music collection servers - for folk like me who already have large music collections - have started to add more sophisticated ways of defining playlists...woo-hoo! :grinning:

## So what is this fork for?

Navidrome has been adding support for smart playlists for sometime now, great..._however...there are things I want a smart playlist to do that are not currently possible_

I wanted to be able to do more...after all, more is more better, right?  
I am an album-centric music listener _(is that a thing?  If it wasn't, it is now)_:rofl: and really like the "Random Album" functionality in the Navidrome GUI.
I want something similar for playlists so I can easily trigger a random album from a client such as Amperfy.  Or play the playlist on my (somewhat antiquated) Sonos system.

So, I've had a quick play around with using SQL as the language to define playlists:

_The music collection is a database so a SQL style smart query syntax should allow flexible, sophisticated queries to be created as playlist definitions_

## You can
 
 - Make complex (_or simple, not every playlist has to be complicated!_) queries using SQL syntax
 - BTW, if you don't know SQL, this could be a good excuse to begin your learning!
 - Your SQL is validated so you can not do anything silly like delete your database tables
 - Use SQL keywords including DISTINCT LIKE LIMIT
 - Use SQL behaviours such as subqueries
 
## Which means you can
 
 Pick a random album
 - Pick a random album by genre
 - Pick a random album by (artist/album artist)
 - Pick a random album by a list of (artists/album artists)
 Write smart queries that include play counts, favourites, etc as criteria...
 _and many more things_ :wink:

## How about some examples?

Sure.  Pick a random album...

```
PLAYLIST name: Random Album, description: Pick an album at random
album_id=(SELECT id FROM album ORDER BY random() LIMIT 1)
ORDER BY disc_number, track_number ASC
```

Pick a random album by Bill Frisell _(other artists are available)_

```
PLAYLIST name: Frisell album, description: Pick a random Bill Frisell album
album_id = (SELECT id FROM album WHERE album_artist = 'Bill Frisell' ORDER BY random() LIMIT 1)
ORDER BY disc_number, track_number ASC
```

Pick 5 random Beatles tracks

```
PLAYLIST name: Beatles x5, description: Pick 5 random Beatles songs
id IN (SELECT id FROM media_file WHERE album_artist = 'The Beatles' ORDER BY random() LIMIT 5)
```

Pick a random album from a list of AC/DC albums

```
PLAYLIST name: AC/DC Bon Scott, description: Pick a random AC/DC with Bon Scott album from a list
album_id =
(SELECT id FROM album WHERE album_artist = 'AC/DC' 
AND name IN ('Let There Be Rock', 'Powerage', 'Live From the Atlantic Studios')
ORDER BY random() LIMIT 1)
ORDER BY disc_number, track_number ASC
```

Or use a date range

```
PLAYLIST name: AC/DC Bon Scott, description: Pick a random AC/DC with Bon Scott album released between 1974 and 1979
album_id =
(SELECT id FROM album WHERE album_artist = 'AC/DC' 
AND min_year >='1974'
AND min_year <= '1979'
ORDER BY random() LIMIT 1)
ORDER BY disc_number, track_number ASC
```

How about a favourite by AC/DC

```
PLAYLIST name: AC/DC faves, description: Pick a random AC/DC album from your favourites
album_id =
(SELECT al.id FROM album al, annotation an 
WHERE al.album_artist = 'AC/DC' 
AND an.item_id = al.id
AND an.starred IS true
ORDER BY random() LIMIT 1)
ORDER BY disc_number, track_number ASC
```

## What do I do with these?

Put them in a file ending with `.smq` in your music library (doesn't matter where, but I put all mine in the same directory at the top of my music library).
Navidrome will pick them up when it scans the music library

## How do you develop & test your queries?

You can easily develop, test and tweak your SQL queries against your database (or a copy of) a SQLite-compatible database manager
- For example, I use the open-sourced [SQLiteStudio](https://sqlitestudio.pl/) [github](https://github.com/pawelsalawa/sqlitestudio)
Replace the `PLAYLIST` line with 

```
SELECT album_id, album, disc_number, track_number, title FROM media_file WHERE 
```

_(adjust the fields returned to taste)_

## Shoutouts

I needed a way to validate the SQL.  Maybe I could have done that using Squirel, but I think not, it just looks like (yet another) SQL builder.  What I wanted was a parser with validation.  And [Tree-sitter](https://tree-sitter.github.io/tree-sitter/) seems to tick all the boxes:
- GitHub uses Tree-sitter to support in-browser symbolic code navigation in Git repositories
- [Go bindings coutesy of Maxim Sukharev aka smacker](https://github.com/smacker/go-tree-sitter) with SQL support

# And finally,

__Any feedback is welcome!__ I created this fork for three reasons
- To learn Go _(Go is so cute!)_
- To learn git, github, etc _(grrrr, I still feel like a total noob, neccessary get better :persevere:)_
- To offer the Navidrome community an idea with a working proof-of-concept that I hope they like
There is always room for improvement and constructive criticism is welcome - obviously it's incredibly unlikely this fork is perfectly coded and bug-free

## Caution! This is a fork not an official release!

**Disclaimer**: This fork may or may not be in an unstable or even broken state during development. 
Please only use if you want to experiment with this code - for exmaple, you have already set up your own Navidrome development environment
If you just want the official Navidrome go get it [here](https://github.com/navidrome/navidrome)
