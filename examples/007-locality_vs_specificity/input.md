Locality is established by creating a new dictionary object

This example shows local instructions winning oevr global ones.  Observe development.app1.db.url.  Even thought the global expression is more specific, the specification of a local wildcard expression overrides the global one.  This is because the most local expression wins a conflict. 

You can see that more specific expressions still win if they have the same locality, in `development.app2.db.url`
