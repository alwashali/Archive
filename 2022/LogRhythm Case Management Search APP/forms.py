from wtforms import (
    StringField,
    PasswordField,
)

from flask_wtf import FlaskForm
from wtforms.validators import InputRequired, Length


class login_form(FlaskForm):
    username = StringField(validators=[InputRequired(), Length(1, 64)])
    pwd = PasswordField(validators=[InputRequired(), Length(min=8, max=72)])
    # Placeholder labels to enable form rendering


