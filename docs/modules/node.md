# Node

A node is a data stored in PostgreSQL engine. A node has a set of properties mainly defined as a general purposes.

-   Nid: The public reference of the node
-   Type: The code used to identify the handler used to manage the node
-   Name: The name, it can be used for internal or external purpose, ideally only for internal.
-   Slug: The slug used for the node, there is not unique constraint on this field, your are free to add this constraint in the validation process or with a PostreSQL index (`CREATE UNIQUE INDEX user_slug ON nodes (col_a) WHERE (type == "core.user");`)
-   Status:
    -   0 - New: The node has been created
    -   1 - Draft: The node is currently being edited
    -   2 - Completed: The node is completed and can be reviewed
    -   3 - Validated: The node is validated and ready for production usage
-   Weight: The weight of the node versus other node, can be used to reorder a list
-   Revision: The current revision of the node, a new revision is created on updated
-   Version: The current data and meta structure version, this can be used to run migrations.
-   CreatedAt
-   UpdatedAt
-   Enabled: The node is enabled and can be used in production.
-   Deleted: The node is deleted and can be removed by a background task.
-   Parents (array of FK Node): A list of parent's node if you use a hierarchical view for this node.
-   UpdatedBy (FK Node): The author of the last update
-   CreatedBy (FK Node): The original author
-   ParentNid (FK Node): The direct parent's node if you use a hierarchical view for this node.
-   SetNid (FK Node): If the node has been created with other nodes, there are all part of the same set.
-   Source (FK Node): If the node has been created from a node, ie a thumbnail from a YouTube's video
-   Data: a structure stores as a JSONb, it hold the user's input or system's input
-   Meta: a structure stores as a JSONb, it hold the related meta from a node
-   Access: an array of string required to manipulate the node.

The current description does not force you about how to use a node for your usage, it is just a guide line. You are free to use the api at your will and you are free to query deleted or un-completed nodes.

## Data vs Meta

The data field is the `Data` coming from your users or api and the `Meta` is more about storing about internal information. Let's take an example with a `User` node.

The user's ``Data` might look like:

    type User struct {
        FirstName   string   `json:"firstname"`
        LastName    string   `json:"lastname"`
        Email       string   `json:"email"`
        Locked      bool     `json:"locked"`
        Enabled     bool     `json:"enabled"`
        Expired     bool     `json:"expired"`
        Username    string   `json:"username"`
        Password    string   `json:"password"`
        NewPassword string   `json:"newpassword,omitempty"`
    }

The data stores all editable information or usable information for an enduser point of view.

Where the user's `Meta` is:

    type UserMeta struct {
        PasswordCost int     `json:"password_cost"`
    }

The Meta stores internal information such the password cost

Again those simple rules are just guideline, you might not want to use the Meta at all. It is up to you!
