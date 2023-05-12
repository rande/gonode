{% extends "prism:layouts/error.tpl" %}

{% block page_title %}Page Not Found{% endblock %}
{% block content %}
    <h1>Page Not Found</h1>

    <div>
        The requested page cannot be found. <br />
        <em>(no prism source found.)</em>
    </div>
{% endblock %}