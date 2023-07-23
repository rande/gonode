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
                <li><a href="{{ url('dashboard_node_edit', url_values('uuid', elm.Uuid)) }}">{{ elm.Name }} - {{ elm.Type }} - {{elm.Uuid }}</a></li>
            {% endfor %}
        </ul>
    </div>
{% endblock %}}