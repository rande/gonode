<form action="POST" action="{{ action }}" encoding="{{ encoding }}">
{% for field in form.Fields %}
    {{ form_field(field, form) }}
{% endfor %}
</form>