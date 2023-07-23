{% extends "layouts/base.tpl" %}

{% block content %}
    <h1>Title: {{ node.Name }}</h1>

    <div>
        <h2>{{ node.Data.Title }}</h2>
        {{ node.Data.Content }}
    </div>
{% endblock %}