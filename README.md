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
Let's start with a vanilla JSON example, `example1.json`

    {
        "development": {
            "app1.db.user" : "app1_user",
            "app1.db.password" : "secret_1",
            "app1.db.url" : "jdbc:h2:/sample/path",
            

            "app2.db.user" : "app2_user",
            "app2.db.password" : "secret_2",
            "app2.db.url" : "jdbc:h2:/sample/path",
            

            "app3.db.user" : "app3_user",
            "app3.db.password" : "secret_3",
            "app3.db.url" : "jdbc:h2:/sample/path"
        }
    }


Let's use our command line tool to extract the actual values

    $ tool --file example1.json development.app1.db.user
    app1_user

## Wildcard matching
The first thing this demonstrates is that our tool can read vanilla JSON.  This is a very useful starting point, because it will allow us to import existing JSON configs directly.  

However, there is a lot of repetition in this config.  For example, you can see that all three database are configured to use the same url.  

    $ tool --file example1.json development.app1.db.url
    jdbc:h2:/sample/path
    $ tool --file example1.json development.app2.db.url
    jdbc:h2:/sample/path
    $ tool --file example1.json development.app3.db.url
    jdbc:h2:/sample/path

Let's use a wildcard to remove this repetition.  This is represented by the `*` charecter

    {
        "development": {
            "app1.db.user" : "app1_user",
            "app1.db.password" : "secret_1",
            
            "app2.db.user" : "app2_user",
            "app2.db.password" : "secret_2",
            
            "app3.db.user" : "app3_user",
            "app3.db.password" : "secret_3",
            
            "*.db.url" : "jdbc:h2:/sample/path"
         }     
    }

Let's query the app1 url to see the wildcard in action.

    $ tool --file example1.json localhost.app1.db.url
    jdbc:h2:/sample/path

You can see that the wildcard value is being used.

## Specifying a convention
You can see a convention in this example as well.  Each of the users is named after the db they are connecting to.  Let's extract that path into a variable, `app_name`, and substitute it in the returned value.  This is done by appending the variable to the path with a `$`. 

    {
        "development": {
            "app1.db.password" : "secret_1",
            
            "app2.db.password" : "secret_2",
            
            "app3.db.password" : "secret_3",
            
            "*$app_name.db.user" : "${app_name}_user"
            "*.db.url" : "jdbc:h2:/sample/path"
         }     
    }

Let's query the app1 url to see the substitution in action

    $ tool --file example1.json localhost.app1.db.user
    app1_user


# Testing the format
The tests are centered around stdin & stdout, in order to ensure equivalence across multiple implementations.  You can see the values in the `test-data` directory
