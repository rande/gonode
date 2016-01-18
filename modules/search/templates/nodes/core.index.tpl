{% extends "layouts/base.tpl" %}

{% block content %}
    <h1>Name: {{ node.Name }} - {{ node.Type }}</h1>

    <div>
        <ul>
            <li>Page: {{ pager.Page }}</li>
            <li>Per Page: {{ pager.PerPage }}</li>
        </ul>

        <ul>
            {% for elm in pager.Elements %}
                <li><a href="/prism/{{ elm.Uuid }}">{{ elm.Name }}</a> - {{ elm.Type }}</li>
            {% endfor %}
        </ul>

        <ul>
            <li><a href="/prism/{{ node.Uuid }}?page={{ pager.Previous }}">Previous</li>
            <li><a href="/prism/{{ node.Uuid }}?page={{ pager.Next }}">Next</li>
        </ul>
    </div>
{% endblock %}}