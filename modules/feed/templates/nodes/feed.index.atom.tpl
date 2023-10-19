<?xml version="1.0" encoding="utf-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">
    <title>{{ node_data(node, "Title") }}</title>
    <link href="feeds/system.atom" rel="self" />

    {% for elm in pager.Elements %}
        <entry>
            <title>{{ node_data(elm, "Title") }}</title>
            <link rel="alternate" type="text/html" href="{{ url("prism", url_values("nid", elm.Nid), request_context) }}"/>
            <id>{{ elm.Nid }}</id>
            <summary><![CDATA[{{ node_data(elm, "Content") }}]]></summary>
        </entry>
    {% endfor %}
</feed>