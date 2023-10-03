{% if field.Input.Type == 'form' %}
    {% for child in field.Children %}
        {{ form_field(child.Name, field.InitialValue) }}
    {% endfor %}
{% else %}
    {{ form_label(field.Name, form) }}
    {{ form_input(field.Name, form) }}
    {{ form_help(field.Name, form) }}
    {{ form_errors(field.Name, form) }}
{% endif %}