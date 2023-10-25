{%if true %}<select name="{{ input.Name}}" value="{{input.Value}}" id="{{ input.Id }}"{% endif %}
    {% if input.Class %} class="{{ input.Class }}" {% endif %}
    {% if input.Style %} style="{{ input.Style }}" {% endif %}
    {% if input.Required %} required {% endif %}
    {% if input.Readonly %} readonly {% endif %}
    {% if input.Disabled %} disabled {% endif %}
    {% if input.Multiple %} multiple{% endif %}
    {% if input.Autofocus %} autofocus{% endif %}
    {% if input.Novalidate %} novalidate{% endif %}>

    {% for option in field.Children %}
        <option value="{{ option.Input.Value }}" id="{{option.Input.Id}}" {% if option.Input.Checked %} selected {% endif %}>{{ option.Label.Value }}</option>
    {% endfor %}
{%if true %}</select>{% endif %}