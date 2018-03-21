Form
=====


### Email Type

    <Field label="Email:" name={"model.email"} type="email" validators={[Email()]} />
    
### Textarea Type
    
    <Field label="Description:" name={"model.description"} type="textarea" />
    
### Select Type

    <Field label="Select:" name={"model.select"} type="select" >
        <option value="1">One</option>
        <option value="2">Two</option>
    </Field>

### Multiselect Type

    <Field label={"Multiselect:"} name={"model.selectmultiple"} type="select-multiple">
        <option value="thomas">Thomas</option>
        <option value="estelle">Estelle</option>
        <option value="manon">Manon</option>
    </Field>

### Checkbox Type

    <Field type="checkbox" value="one" name={"model.choice.one"} label={"One"} disabled />
    <Field type="checkbox" value="two" name={"model.choice.two"} label={"Two"} />
    <Field type="checkbox" value="three" name={"model.choice.three"} label={"Three"}/>

### Radio Type

    <Field type="radio" value={1} name={"model.radio"} label={"1"}/>
    <Field type="radio" value={2} name={"model.radio"} label={"2"}/>
    <Field type="radio" value={3} name={"model.radio"} label={"3"}/>