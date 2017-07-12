This is an example from real life.  The original file was approximately 360 lines of vanilla JSON.  This replacement version comes in at about 40 lines.  An order of magnitude improvment.  Not only is this a smaller file, but the real gains come when extending your system

* Adding new environments is a breeze, it will only require adding a top level url
* Adding a new database will usually require adding a user entry to the `*` locality

This drastically cuts down on the amount of busywork that is required for configuration
