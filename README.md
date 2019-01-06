# NetQuery

## Intro

NetQuery is a project designed to make intrusion detection systems (IDS) available
 to network defenders and security engineers as native storage devices. In this case,
  to make IDSs look like SQL databases.

## Why

Today, the state of network security is pretty dismal.

All of the tools available are generally designed to spew out lots of structured
 and dense log file. This used to be great; but, in todays world this no longer cuts it.

Generally, the solution has been to talk *all of the data* and stream/batch it
 back to a piece of infrastructure build to store a large quantity of data in the
  hopes that some brave soul will figure out the magic query to make it all work.

But thats lame.

## Small Data Over Big Data

NetQuery is designed with a focus on small data. Specifically with IDSs, this means
 being able to get meaningful data out of a single IDS and finding patterns and norms without
needing to store terabytes of logs in the hopes of finding the needle in the needle stack.
