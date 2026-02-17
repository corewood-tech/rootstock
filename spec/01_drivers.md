# ROOTSTOCK by Corewood

## Requirements Specification

> Following the Volere Requirements Template

---

# Section 1: Project Drivers

## 1.1 The Purpose of the Project

Rootstock is an open-source scientific data collection platform that connects researchers with citizen scientists through consumer IoT devices.

### Background

Current methods of acquiring data for scientific analysis are inefficient, expensive, and often insufficient in scope. Meanwhile, consumer electronics — smartphones, weather stations, air quality monitors, water sensors, GPS-enabled devices — represent distributed "mini-labs" already deployed at massive scale. The **last-mile problem** remains: there is no standard, low-friction way to connect these consumer devices to active research needs.

Given the ubiquity of suitable data collection equipment and the critical importance of accurate, reliable data to scientific research, **Project Rootstock** aims to bridge the gap between researchers who need field data and the citizen scientists equipped to provide it.

### Motivation

Enable researchers to request and acquire high-quality field data through an activated, engaged network of citizen scientists — at a fraction of the cost of traditional data collection methods.

---

## 1.2 Goals of the Project

| Goal | Description |
|------|-------------|
| **Research Data Requests** | Researchers define structured data campaigns specifying what data they need, where, when, and to what quality standard. |
| **Citizen Science Response** | Citizen scientists discover campaigns relevant to their location and equipment, connect their IoT devices, and contribute data. |
| **Reward & Recognition** | Gamification and sweepstakes-based incentives keep citizen scientists engaged without requiring per-unit compensation. |

### Success Measurement

The platform succeeds when:

- **Data volume** increases relative to traditional collection methods for equivalent research scopes.
- **Data quality** meets or exceeds researcher-defined thresholds.
- **Relative cost of data acquisition** drops significantly compared to traditional field collection.
- **Community growth** is organic — rewarded citizen scientists recruit others through word-of-mouth.

### Cost Model

Costs are kept low by offering sweepstakes entries, experiences, and recognition to data collectors rather than compensating all collectors based on time or volume. This model scales: as more citizen scientists join, per-unit data cost approaches zero while coverage expands.

---

## 1.3 Stakeholders

### The Client

| Stakeholder | Role | Description |
|-------------|------|-------------|
| **Corewood** | Builder / Platform Owner | Develops and maintains the Rootstock platform. Sets technical direction and platform governance as an open-source project. |

### The Customers

| Stakeholder | Role | Description |
|-------------|------|-------------|
| **Research Institutions** | Data Consumer (Organization) | The institution at which a researcher operates. May be a university, private company, government agency, or independent research organization. Institutions may have compliance, ethics board, or data governance requirements that constrain how data is collected and used. |
| **Researchers** | Data Consumer (Individual) | Individuals who define data needs and analyze collected data on behalf of a research institution. Primary users of the campaign management interface. |
| **Grant Writers / Boards / Oversight Bodies** | Governance | Political, financial, or ethical oversight of research institutions and/or individual researchers. May influence what data can be requested, how it is stored, and what consent is required from participants. |

### Other Stakeholders

| Stakeholder | Role | Description |
|-------------|------|-------------|
| **Citizen Scientists** | Data Producer | Anyone willing to donate time, effort, and physical exertion to further scientific research. They connect personal IoT devices to active campaigns and contribute sensor readings. Their engagement is voluntary and incentivized through gamification rather than direct compensation. |

---

## 1.4 Hands-On Product Users

These are the stakeholders who directly interact with the platform day-to-day.

| User | Goal | Key Need |
|------|------|----------|
| **Researcher** | Gather data to support or refute hypotheses | Define campaigns precisely, receive data that meets quality thresholds, export/analyze results with confidence in provenance |
| **Citizen Scientist** | Contribute to science, be recognized for contributions | Low-friction device onboarding, clear feedback that their data matters, visible recognition and reward |

### Non-Human Actors: IoT Devices

IoT devices are not stakeholders, but they are a **primary source of volatility** that must be acknowledged early.

- **Hardware fragmentation**: Sensors vary wildly in precision, calibration, sampling rate, and communication protocol across manufacturers and models.
- **Firmware instability**: Devices receive OTA updates outside our control that may change data formats, sampling behavior, or connectivity patterns.
- **Connectivity unreliability**: Devices operate in the field — intermittent connectivity, power loss, and environmental interference are the norm, not the exception.
- **Data trust boundary**: Every reading from an IoT device crosses a trust boundary. The platform cannot assume readings are accurate, timely, or unmanipulated.

This volatility directly shapes architecture decisions: the platform must treat device data as untrusted input, validate against campaign-defined quality constraints, and degrade gracefully when devices behave unexpectedly.

---

## 1.5 Stakeholder Priorities

> To be defined in collaboration with stakeholders. Initial assumptions:

| Stakeholder | Primary Concern | Risk if Unaddressed |
|-------------|----------------|---------------------|
| Research Institutions | Data provenance, quality assurance, regulatory compliance | Platform produces data that cannot be used in published research |
| Researchers | Low-friction campaign creation, data accessibility, sufficient geographic/temporal coverage | Platform is too complex to use or doesn't attract enough contributors in needed areas |
| Oversight Bodies | Ethical data collection, participant consent, data privacy | Platform creates liability for the institution |
| Citizen Scientists | Easy device onboarding, clear contribution value, fair recognition | Contributors churn due to friction or feeling unappreciated |
| Corewood | Sustainable open-source model, platform adoption, technical quality | Project stalls or fails to reach critical mass |

---

*Next: [Section 2 — Project Constraints](./02-project-constraints.md)*
