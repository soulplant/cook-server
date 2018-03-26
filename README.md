Cook Database Design Doc

Overview
There are several design concerns that need to be addressed when representing recipes in a relational database.


Relational Layout of a Recipe
When we have something like a recipe that has two lists of things (i.e. instructions and ingredients), we ultimately need to do two queries. The general rule being that a query can be for one set of similar things. This is ultimately what SQL gets us - each query has the same table structure, so aggregates are going to inherently require multiple queries to construct. Or, put another way, one query equals one list of flat objects. Therefore we shouldn't expect to be able to construct a fully denormalised recipe from a single query. Given a single recipe, we can construct all the details needed to render a page in three queries: one for the recipe, and one each for the list of ingredients and steps. That's the fully denormalised approach. I could imagine the steps being inlined in a JSON blob, because they don't really need identity, but it's not really an issue.

So what if we wanted to query a list of recipes? We could do a batch mode. Query the list of recipes you want. Then take those ids and query the other things using a "recipe_id IN (...)" query (and likewise for instructions). We'd then do something like keep a map of the recipes in memory and then iterate through the other queries associating ingredients and instructions to the appropriate recipe incrementally before returning them all to the caller.

How this looks in an ORM like gorm seems fairly straightforward. There'll be some mechanism for joins and then "hydration" will be done in a purely programmatic way using good ol' for loops.


Denormalising Recipes
Each list associated with a recipe implies an additional query for reading the aggregate. These additional queries may need joins of their own and so the complexity of the query is not quite trivial, though certainly it's still fairly easy to grasp. It seems likely that there's an opportunity to reduce the cognitive burden by storing some of these lists inline in the recipe. However, it's not just about removing the relational aspect of those lists, but also denormalising the entities involved. It's possible to use functions that synthesise the relational tables out of JSON arrays (i.e. json_array_values() is a set returning function (SRF) which can be used to turn

  x | [a,b,c]

  into

  x | a
  x | b
  x | c

  (see [1])

But if you're just going to convert the JSON to a relation anyway, then why not just start with relations? It saves the overhead of having to pick apart JSON and you get a nice schema for your data as well.

The only possibility that's interesting is if you can have a query that turns something like:

  x | {ingredients: [1,2]}

  into

  x | {ingredients: [{id: 1, name: "foo"}, {id: 2, name: "bar"}]

I think it may be possible to achieve something like the above using fancy queries and JSON creation functions (there must be a way to create a JSON value from a table), but that seems like more trouble than its worth for our purposes - to make that performant one would want indexes over the JSON fields to be defined, and also writing such explicitly psql SQL harms portability and simplicity.

GORM supports preloading [2] which makes these kinds of queries quite effortless to write, so that's what we shall use.

[1] https://www.periscopedata.com/blog/the-lazy-analysts-guide-to-postgres-json
[2] http://doc.gorm.io/crud.html#preloading-eager-loading
