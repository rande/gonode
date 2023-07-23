{% extends "dashboard:layouts/default.tpl" %}

{% block content %}
    <h1>Search</h1>

    <div>
        <ul>
            <li>Page: {{ pager.Page }}</li>
            <li>Per Page: {{ pager.PerPage }}</li>
        </ul>

        <ul>
            {% for elm in pager.Elements %}
                <li><a href="{{ prism_path(elm) }}">{{ elm.Name }}</a> - {{ elm.Type }}</li>
            {% endfor %}
        </ul>

        <ul>
            <li><a href="{{ prism_path(node, url_values("page", pager.Previous)) }}">Previous</li>
            <li><a href="{{ prism_path(node, url_values("page", pager.Next)) }}">Next</li>
        </ul>
    </div>
{% endblock %}}