Managing configuration is hard.  One of the main risks is that there is a lot of needless repetition in a configuration file, which introduces noise.  This project is an effort to increase the amount of signal in a configuration and reduce the noise.  Some core precepts:

* Configurations always get values as a path
* Configirations are heirarchical
* Configurations are read centric
* It should be possible to place convention on top of configuration

# Goals
The goals of this project are:

* Provide human friendly way to concisely define a config
* Provide a human friendly explanation as to why as certain value was chosen
* Support multiple platforms (java, ruby, go, js, .NET)
* Support vanilla JSON as an config format
* Allow serialized config information to be tranmitted easily

# The format
Please see examples/README.md for more details

# Testing the format
The tests are centered around stdin & stdout, in order to ensure equivalence across multiple implementations.  You can see the values in the `examples` directory
