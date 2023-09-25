{%if true %}<label for="{{ input.Id }}">{{ label.Value }}{% endif %}
    {% if label.Class %}class="{{ label.Class }}" {% endif %}
    {% if style.Style %}class="{{ style.Style }}" {% endif %}
</label>