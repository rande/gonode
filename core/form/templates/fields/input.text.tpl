{%if true %}<input name="{{ input.Name}}" type="{{ input.Type }}" value="{{input.Value}}" id="{{ input.Id }}"{% endif %}
    {% if input.Class %}class="{{ input.Class }}" {% endif %}
    {% if input.Style %}class="{{ input.Style }}" {% endif %}
    {% if input.Placeholder %} placeholder="{{ input.Placeholder }}"{% endif %}
    {% if input.Required %} required {% endif %}
    {% if input.Readonly %} readonly {% endif %}
    {% if input.Disabled %} disabled {% endif %}
    {% if input.Checked %} checked {% endif %}
    {% if input.Multiple %} multiple{% endif %}
    {% if input.Autofocus %} autofocus{% endif %}
    {% if input.Novalidate %} novalidate{% endif %}
    {% if input.Size > 0 %} size="{{ input.Size }}"{% endif %}
    {% if input.MaxLength != 0 || input.MinLength != 0 %} maxlength="{{ input.MaxLength }}"{% endif %}
    {% if input.MaxLength != 0 || input.MinLength != 0 %} minlength="{{ input.MinLength }}"{% endif %}
    {% if input.Min != 0 || input.Max != 0 %} max="{{ input.Max }}"{% endif %}
    {% if input.Min != 0 || input.Max != 0 %} min="{{ input.Min }}"{% endif %}
    {% if input.Step > 0 %} step="{{ input.Step }}"{% endif %}
    {% if input.Height > 0 %} height="{{ input.Height }}"{% endif %}
    {% if input.Width > 0 %} width="{{ input.Width }}"{% endif %}
    {% if input.Pattern %} pattern="{{ input.Pattern }}" {% endif %}
    {% if input.Autocomplete %} autocomplete="{{ input.Autocomplete }}"{% endif %}
    {% if input.List %} list="{{ input.List }}"{% endif %}
>