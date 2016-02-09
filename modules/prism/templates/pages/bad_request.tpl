{% extends "layouts/error.tpl" %}

{% block page_title %}Bad Request{% endblock %}
{% block content %}
    <h1>Bad Request</h1>

    <div>
        The request is not valid.<br />
        <em>(no prism handler found.)</em>
    </div>
{% endblock %}