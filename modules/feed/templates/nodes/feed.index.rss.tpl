<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
    <channel>
        <title>{{ node_data(node, "Title") }}</title>
        <link>{{ request_context.Prefix }}{{ request.URL }}</link>
        <description>{{ node_data(node, "Description") }}</description>

        {% for elm in pager.Elements %}
            <item>
                 <title>{{ node_data(elm, "Title") }}</title>
                 <link>{{ url("prism", url_values("nid", elm.Nid), request_context) }}</link>
                 <description><![CDATA[{{ node_data(elm, "Content") }}]]></description>
                 <pubDate>{{ node_data(elm, "PublicationDate") }}</pubDate>
                 <gui>{{ elm.Nid }}</gui>
            </item>
        {% endfor %}
    </channel>
</rss>