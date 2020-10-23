import { Formik } from "formik";
import React from "react";
import { v4 as uuidv4 } from 'uuid';

const Join = (props) => {
    const initialValues = {
        username: "",
        id: uuidv4(),
        buyin: 100,
    }

    const onSubmit = (values, { setSubmitting }) => {
        setSubmitting(false);
        props.onJoin(values);
    }

    return (
        <div>
            <h1>Join Form</h1>
            <Formik
                initialValues={initialValues}
                onSubmit={onSubmit}
            >
                {({
                    values,
                    errors,
                    handleChange,
                    handleSubmit, 
                    isSubmitting
                }) => (
                  <form onSubmit={handleSubmit}>
                      <input type="text" name="username" onChange={handleChange} value={values.username}/>
                      <input type="text" name="id" onChange={handleChange} value={values.id}/>
                      <input type="number" name="buyin"onChange={handleChange} value={values.buyin}/>
                      <input type="submit" disabled={isSubmitting} value="Join" />
                  </form>  
                )}
            </Formik>
        </div>
    );
}

export default Join;